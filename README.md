This is an app designed to upload files to a backup device with an easy to use UI. It's still on development, so maybe check it out later to see the current progress :)

## Requirements

- [go](https://go.dev/) `1.23.2` or higher
- [bun](https://bun.sh/) `1.1.22` or higher
- [docker](https://docs.docker.com/engine/install/)

## How to use it

1. Clone the repository

```bash
git clone git@github.com:julio-salas03/local-file-vault.git
```

2. Create a `.env` file

```bash
cp .env.example .env
```

3. Run the `init` command

```bash
bun run init
```

**That's it!** The application is now running on port `8080` ready to backup your files!

### Optional for development

- [air](https://github.com/air-verse/air) `1.61.1` or higher

## Notes for development

The `npm run dev` command with span 2 servers with `concurrently`. You'll want to use the one running on port `3000`, as that's the vite server and comes with HMR. Avoid using the server on port `8080` unless you're testing api related logic
