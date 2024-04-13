![ReadMe Banner](https://github.com/MathisVerstrepen/github-visual-assets/blob/f6c2cfde9b430e98f80bbc3d11993117090400bb/banner/Letterboxd-Jellyfin-Go.png?raw=true)

# Letterboxd Jellyfin Go Integration

This is a simple script that allows you to import your Letterboxd watchlist into Jellyfin. 

It will scan your watchlist and add any movies that are not already in your Jellyfin library via Radarr. It will also manage a watchlist collection in Jellyfin that will be updated with the movies that are in your watchlist.

This script is a rewrite in Go of the original script that was written in Python. The original script can be found [here](https://github.com/MathisVerstrepen/letterboxd-jellyfin).

![Splitter-1](https://raw.githubusercontent.com/MathisVerstrepen/github-visual-assets/main/splitter/splitter-1.png)


## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Deploying](#deploying)
- [License](#license)

![Splitter-1](https://raw.githubusercontent.com/MathisVerstrepen/github-visual-assets/main/splitter/splitter-1.png)


## Features

- [x] Scan Letterboxd watchlist by scraping the website with self made Go scraper.
- [x] Add movies to Jellyfin library via Radarr API.
- [x] Manage a watchlist collection in Jellyfin that will be updated with the movies that are in your watchlist and remove the movies that you have watched.
- [ ] Handle Mini-Series on Letterboxd and import them into Sonarr.

![Splitter-1](https://raw.githubusercontent.com/MathisVerstrepen/github-visual-assets/main/splitter/splitter-1.png)


## Installation

**Warning:** This project is still in development and lots of features are hard coded. This means that you will have to change the code to make it work for your specific setup.

1. Clone the repository:

    ```shell
    git clone https://github.com/MathisVerstrepen/letterboxd-jellyfin-go
    ```

2. Build the project:

    ```shell
    go build
    ```

3. Run the project:

    ```shell
    ./letterboxd-jellyfin-go
    ```
![Splitter-1](https://raw.githubusercontent.com/MathisVerstrepen/github-visual-assets/main/splitter/splitter-1.png)


## Deploying

This script is not meant to be run on another system than mine. It is not very user friendly and I have no intention of making it so. If you want to use it, you will have to modify it to suit your needs.

Currently, the script is deployed as a docker container and run as a cronjob every 2 minutes. The docker container is built automatically using my custom deployment pipeline that can be found [here](https://github.com/MathisVerstrepen/ApolloLaunchCore).

![Splitter-1](https://raw.githubusercontent.com/MathisVerstrepen/github-visual-assets/main/splitter/splitter-1.png)


## License

This project is licensed under the [MIT License](LICENSE).