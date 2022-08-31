# rofi-gitlab

> **Important** Work in Progress

A rofi mode for fast gitlab navigation

## Installation

Make sure you have rofi installed:

```
sudo pacman -S rofi
```

Run `go build` and `go install`. Make sure you have `$HOME/go/bin` in your `$PATH`!

Create a `config.json` at `$HOME/.config/rofi-gitlab` and change the values accordingly:

```
{
 "BaseUrl": "https://gitlab.example.com",
 "Token": "<your private token>",
 "TTL": 3600, 
}%
```

TTL is the number of seconds when `rofi-gitlab` will cache the results of the projects.

## Usage 

Use rofi-gitlab as mode:

```
$ rofi -modi "gitlab:$HOME/go/bin/rofi-gitlab" -show gitlab
```

