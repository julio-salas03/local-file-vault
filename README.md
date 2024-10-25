This is an app designed to upload files to a backup device with an easy to use UI. It's still on development, so maybe check it out later to see the current progress :)

## Requirements

- [go](https://go.dev/) `1.23.2` or higher
- [bun](https://bun.sh/) `1.1.22` or higher

### Optional for development

- [air](https://github.com/air-verse/air) `1.61.1` or higher

## Notes for development

The `npm run dev` command with span 2 servers with `concurrently`. You'll want to use the one running on port `3000`, as that's the vite server, which comes with HMR. Avoid using the server on port `8080` unless you're testing api related logic

## Troubleshooting

### Cannot run `air` when running `npm run dev`

Make sure you have the following alias in your bash aliases declarations `alias air='$(go env GOPATH)/bin/air'`
