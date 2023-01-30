package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"localcast/models"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

//go:embed assets templates
var embeddedFiles embed.FS

func podcastCover(context *gin.Context) {
	// covert podcast id to int
	podcastID, err := strconv.Atoi(context.Param("podcast"))
	if err != nil {
		context.String(http.StatusNotFound, "Podcast not found")
		return
	}

	// get podcast
	var podcast models.Podcast
	models.DB.Where("id = ?", podcastID).First(&podcast)

	// get path from router
	path, exists := context.Get("path")

	if !exists {
		context.String(http.StatusNotFound, "Location not found")
		return
	}


	podcastPath := fmt.Sprintf("%s/Downloads/%s", path, podcast.DownloadFolder)

	// check if podcastPath/folder.jpg or podcastPath/folder.png exists
	// if so, serve that
	if _, err := os.Stat(fmt.Sprintf("%s/folder.jpg", podcastPath)); err == nil {
		context.File(fmt.Sprintf("%s/folder.jpg", podcastPath))
		return
	} else if _, err := os.Stat(fmt.Sprintf("%s/folder.png", podcastPath)); err == nil {
		context.File(fmt.Sprintf("%s/folder.png", podcastPath))
		return
	}
}

func episodeAudio(context *gin.Context) {
	// covert eisode id to int
	episodeID, err := strconv.Atoi(context.Param("episode"))
	if err != nil {
		context.String(http.StatusNotFound, "Episode not found")
		return
	}

	// get episode
	var episode models.Episode
	models.DB.Where("id = ?", episodeID).First(&episode)

	// check if episode exists
	if episode.ID == 0 {
		context.String(http.StatusNotFound, "Episode not found")
		return
	}

	// get podcast
	var podcast models.Podcast
	models.DB.Where("id = ?", episode.PodcastId).First(&podcast)

	// check if podcast exists
	if podcast.ID == 0 {
		context.String(http.StatusNotFound, "Podcast not found")
		return
	}

	// get path from router
	path, exists := context.Get("path")

	if !exists {
		context.String(http.StatusNotFound, "Location not found")
		return
	}

	// get episode path
	episodePath := fmt.Sprintf("%s/Downloads/%s/%s", path, podcast.DownloadFolder, episode.FileName)

	context.File(episodePath)
}

func landingPage(context *gin.Context) {
	// get all podcasts
	var podcasts []models.Podcast

	// get podcasts where there are episodes with FileName != ""
	models.DB.Where("id IN (SELECT podcast_id FROM episode WHERE download_filename != '')").Find(&podcasts)

	context.HTML(http.StatusOK, "index.html", gin.H{
		"podcasts": podcasts,
	})
}

func podcastPage(context *gin.Context) {
	// covert podcast id to int
	podcastID, err := strconv.Atoi(context.Param("podcast"))
	if err != nil {
		context.String(http.StatusNotFound, "Podcast not found")
		return
	}

	// get podcast
	var podcast models.Podcast
	models.DB.Where("id = ?", podcastID).First(&podcast)

	if podcast.ID == 0 {
		context.String(http.StatusNotFound, "Podcast not found")
		return
	}

	// convert podcast description to html
	description := template.HTML(blackfriday.MarkdownCommon([]byte(podcast.Description)))
	description = template.HTML(bluemonday.UGCPolicy().SanitizeBytes([]byte(description)))

	// create a short description
	var shortDescriptionRaw string
	var moreStatus bool

	if len(podcast.Description) < 400 {
		moreStatus = false
		shortDescriptionRaw = podcast.Description
	} else {
		moreStatus = true
		shortDescriptionRaw = podcast.Description[:strings.LastIndex(podcast.Description[:400], " ")]
	}

	shortDescription := template.HTML(blackfriday.MarkdownCommon([]byte(shortDescriptionRaw)))
	shortDescription = template.HTML(bluemonday.UGCPolicy().SanitizeBytes([]byte(shortDescription)))

	// all episodes with FileName != "" and ordered by date
	var episodes []models.Episode
	models.DB.Where("podcast_id = ? AND download_filename != ''", podcastID).Order("published ASC").Find(&episodes)

	// episode count
	var episodeCount = len(episodes)

	for i := range episodes {
		// set episode description to html
		episodeDesc := template.HTML(blackfriday.MarkdownCommon([]byte(episodes[i].Description)))
		episodes[i].DescriptionHtml = template.HTML(bluemonday.UGCPolicy().SanitizeBytes([]byte(episodeDesc)))

		// set episode date string (21st May 2021, 01:00 AM)
		episodes[i].DateString = time.Unix(int64(episodes[i].Date), 0).Format("2 Jan 2006, 03:04 PM")
	}

	lastUpdated := episodes[len(episodes)-1].DateString
	lastUpdated = lastUpdated[:strings.LastIndex(lastUpdated, ",")]

	context.HTML(http.StatusOK, "podcast.html", gin.H{
		"podcast":          podcast,
		"description":      description,
		"shortDescription": shortDescription,
		"moreStatus":       moreStatus,
		"episodes":         episodes,
		"episodeCount":     episodeCount,
		"lastUpdated":      lastUpdated,
	})
}

func FaviconFS() http.FileSystem {
	sub, err := fs.Sub(embeddedFiles, "./assets/favicon.ico")
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}

func setupRouter(path string) *gin.Engine {
	router := gin.Default()
	
	// set path to router
	router.Use(func(context *gin.Context) {
		context.Set("path", path)
		context.Next()
	})

	router.SetTrustedProxies(nil)

	templ := template.Must(template.New("").ParseFS(embeddedFiles, "templates/*"))
	router.SetHTMLTemplate(templ)

	router.StaticFS("/public", http.FS(embeddedFiles))


	router.GET("/", landingPage)
	router.GET("/podcast/:podcast", podcastPage)
	router.GET("/podcast/:podcast/cover", podcastCover)
	router.GET("/episode/:episode/audio", episodeAudio)

	// server favicon under root
	router.GET("/favicon.ico", func(context *gin.Context) {
		context.FileFromFS(".", FaviconFS())
	})

	fmt.Println("Server running on localhost:3000")

	return router
}

func main() {
	var path = flag.String("path", "", "gPodder directory path")
	var port = flag.String("port", "3000", "port to run server on")
	flag.Parse()

	if *path == "" {
		fmt.Println("Please specify the path to gPodder directory")
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)

	// connect to database
	models.ConnectDatabase(fmt.Sprint(*path, "/Database"))

	router := setupRouter(*path)

	router.Run(":" + *port)
}
