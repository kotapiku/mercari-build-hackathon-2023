# Mini Mercari Web App

## Requirements

* [go](https://go.dev/)

## Getting started

### 1. Update environment values

1. Create `initialize` branch.

2. Run bellow command.

```shell
$ go run tools/setup.go -g [your github name] -t [your team id]

(e.g.) $ go run tools/setup.go -g yourname -t 16
```

3. Create a PR form `initialize` to `main`.

### 2. Launch services

See `backend/README.md` for backend service and `frontend/simple-mercari-web/README.md` for frontend service.

## How to run bench marker(TBD)
Release on Monday. Please wait!

## What should we do first?

- First, stand up services and see logs both of backend and frontend services
  - For backend, you can see the logs on the terminal where the server is set up
  - For frontend, use Chrome Devtool to check
    - https://developer.chrome.com/docs/devtools/overview/
- Try to use mini Mercari and find problems that should not occur in the original Mercari service. For example:
  - When you check the item detail page of your listed items...?
  - When you try to buy items that exceed your available balance...?
  - When there are multiple users purchase an item at the same time...?
- The UI is quite simple and difficult to use
  - It looks inconvenient if there is no message indicating a request to the backend has failed
  - As the number of items increase, the UI is likely to become slow
  - Feel free to make the site more user-friendly like implementing features not implemented in actual Mercari web site

However, your changes must be made within the constraints of the bench marker. Please refer to backend/README.md for details.