# 😡 who-the-hell-unfollowed-me ![Go](https://github.com/wuhan005/who-the-hell-unfollowed-me/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/who-the-hell-unfollowed-me)](https://goreportcard.com/report/github.com/wuhan005/who-the-hell-unfollowed-me) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/wuhan005/who-the-hell-unfollowed-me)

## TL;DR

This is a tool to help you find out who unfollowed you on GitHub.

## Setup

1. Create a [personal access token](https://github.com/settings/personal-access-tokens/new) with `Followers - Read
   only` and `Gists - Read and write` permission.
2. Create a [GitHub Gist](https://gist.github.com/), and save the Gist ID in the last part of the URL.
3. Fork this repository, set `GH_TOKEN` and `GIST_ID` secrets in your forked repository's settings.
4. Wait for the first run to complete, and then you can see the result in your Gist.

## Demo

https://gist.github.com/wuhan005/c954d5b61ed1eb7d15a13b0d80ba0dd8

## License

MIT License
