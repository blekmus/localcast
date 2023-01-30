<p align="center">
  <h1 align="center"><img width="400" src="./assets/logo.png" alt="localcast"></h1>
  <p align="center">A Go powered frontend for <a href="https://github.com/gpodder/gpodder" target="_blank">gPodder</a> podcast downloads</p>
</p>

<br>
<br>

## About

Localcast checks locally downloaded podcast files from gPodder and serves them on the browser. This makes it easier to browse and listen them. 

The server is powered by Gin with SSR html pages. The data is extracted from gPodder's sqlite database. No network requests are made site-wide because this is meant to index archived/downloaded files only.

## Previews

![image](https://user-images.githubusercontent.com/47277246/215434469-338b5533-85d7-4e3d-82be-a230113f004f.png)

![image2](https://user-images.githubusercontent.com/47277246/215437436-61f51406-d47d-4beb-8ec2-f2c897c19fef.png)
