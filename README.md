# Project Setup and Build Guide

This repository uses [Go-Task](https://taskfile.dev/) as its build automation tool to simplify development workflows. Follow the instructions below to configure your environment and run the project build tasks.

## Installation & Setup

```bash
task export-model
```
## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
