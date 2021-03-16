# Developing go-redmine

## Building

### Tools

Building requires these installed tools. In general, this project does not rely on cutting-edge or experimental versions. Decent versions are okay, but if building fails you should check if you have quite outdated versions:

- Make
   - f. e. GNU Make 4.2.1
- Docker
   - f. e. Client/server version 19.03.5
- Golang compiler
   - f. e. 1.14.12

### Local Building and other interesting `make` targets

make target | action
------------|-------
`make` / `make compile` | builds the provider binary in Go container
