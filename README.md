# URL Shortener

A simple URL shortener service, currently inteded for low volume usage.

*N.B. This is a project used for me to learn and experiment with Go, so shouldn't be taken seriously.*

## 🚀 Usage

### 1. Clone the repo

```
git clone https://github.com/tomvanner/url-shortener.git
cd url-shortener
```

### 2. Build the project

Create the necessary sqlite database and build the relevant docker images

```
make build
```

### 3. Start the API

Start the docker container

```
make up
```
