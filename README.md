# README

## About

launch.games Launcher is an app that allows you to launch games on your desktop.
It supports multiple desktop platforms and multiple games.
Launcher provides self-update and game update features.

Launcher is built using the Wails framework.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

When building, you should pass in
the `-ldflags "-X games.launch.launcher/config.LauncherId=put-the-launcher-uuidv4-here"` flag to set the
LauncherId. This is used to identify the launcher and its games.

Tags can be used to build for specific configurations. The following tags are
available: `Development`, `Test`, `Shipping`.

## Database Metadata

- Metadata is requested from the API corresponding to the build configuration (Development, Test, Shipping).
- Required metadata is stored in the database as a record of `Entity` of type `launcher-v2` and
  corresponding `LauncherV2`
  record.
- The `LauncherV2` record has a list of `ReleaseV2` metadata, which is used to determine launcher releases (and
  versions).
- `LauncherV2` also has a list of `AppV2` records, which is used to determine games supported by the launcher.
- Each `AppV2` record has a list of `ReleaseV2` records, which is used to determine the versions of the game.
- `AppV2` record can have a link to the `SdkV2` record, which is used to determine the SDK used by the game.
- `SdkV2` record has a list of `ReleaseV2` records, which is used to determine the versions of the SDK used by the app.

Example of the Launcher metadata stored in the database:

```json
{
  "id": "AAAAAAAA-AAAA-4AAA-AAAA-AAAAAAAAAAAA",
  "createdAt": "2022-12-31T15:14:52.83374Z",
  "owner": {
    "id": "88888888-8888-4888-8888-888888888888",
    "createdAt": "0001-01-01T00:00:00Z",
    "name": "Hackerman"
  },
  "files": {
    "entities": [
      {
        "id": "99999999-9999-4999-9999-999999999999",
        "createdAt": "2022-12-31T15:25:33.522303Z",
        "type": "image-app-icon",
        "url": "https://filesamples.com/samples/image/ico/sample_640%C3%97426.ico",
        "mime": "image/x-icon",
        "platform": "Win64",
        "size": 175370,
        "originalPath": "sample-clouds-400x300.jpg"
      },
      {
        "id": "AAAAAAAA-AAAA-4AAA-AAAA-AAAAAAAAAAAA",
        "createdAt": "2022-12-31T15:25:33.522303Z",
        "type": "image-app-icon",
        "url": "https://download.samplelib.com/png/sample-red-400x300.png",
        "mime": "image/png",
        "size": 1000,
        "originalPath": "sample-red-400x300.png"
      }
    ]
  },
  "name": "Genesis",
  "description": "",
  "releases": {
    "entities": [
      {
        "id": "BBBBBBBB-BBBB-4BBB-BBBB-BBBBBBBBBBBB",
        "createdAt": "2023-01-01T01:59:41.438148Z",
        "files": {
          "entities": [
            {
              "id": "00000000-0000-4000-A000-000000000000",
              "createdAt": "2023-01-01T02:01:11.661659Z",
              "type": "launcher",
              "url": "https://xxxx.xxxx.amazonaws.com/BBBBBBBB-BBBB-4BBB-BBBB-BBBBBBBBBBBB/launcher.exe",
              "mime": "application/vnd.microsoft.portable-executable",
              "size": 11560448,
              "originalPath": "launcher.exe"
            }
          ]
        },
        "version": "1.0.0",
        "codeVersion": "1.0.0",
        "contentVersion": "1.0.0",
        "name": "Genesis",
        "description": "Test Launcher Release 1.0.0",
        "archive": false
      }
    ]
  }
}
```

Example of Launcher `AppV2` metadata list stored in the database:

```json
[
  {
    "id": "11111111-1111-4111-A111-111111111111",
    "createdAt": "2023-01-01T11:06:19.521038Z",
    "owner": {
      "id": "AAAAAAAA-AAAA-4AAA-AAAA-AAAAAAAAAAAA",
      "createdAt": "0001-01-01T00:00:00Z",
      "name": "Hackerman"
    },
    "files": {
      "entities": [
        {
          "id": "21494b48-c5c4-4b0e-b039-3454fa94e386",
          "createdAt": "2023-01-01T11:08:34.022464Z",
          "type": "launcher-app-card",
          "url": "https://samplelib.com/lib/preview/jpeg/sample-clouds-400x300.jpg",
          "mime": "image/jpeg",
          "size": 47419,
          "originalPath": "sample-clouds-400x300.jpg"
        },
        {
          "id": "aba345d8-de86-4f8c-bfe0-e5a69dd7e7e7",
          "createdAt": "2023-01-01T11:08:34.022464Z",
          "type": "launcher-app-background",
          "url": "https://samplelib.com/lib/preview/jpeg/sample-clouds-400x300.jpg",
          "mime": "image/jpeg",
          "size": 47419,
          "originalPath": "sample-clouds-400x300.jpg"
        }
      ],
      "total": 2
    },
    "links": {
      "entities": [
        {
          "id": "438a45ce-7a2d-4815-b020-dabb7bfaa05d",
          "url": "YouTube",
          "name": "https://youtu.be/y2ECgOhoDGs"
        }
      ],
      "total": 1
    },
    "name": "LE7EL",
    "external": false,
    "sdk": {
      "id": "25e66de8-c6b1-4c43-870e-9c03bbe32d61",
      "createdAt": "2023-01-04T10:44:14.699665Z",
      "releases": {
        "entities": [
          {
            "id": "12ffbe75-7831-4d02-b679-ef27a6b279ba",
            "createdAt": "2023-01-04T19:23:20.262856Z",
            "files": {
              "entities": [
                {
                  "id": "8f268da6-8636-4349-9994-872cd5512fac",
                  "createdAt": "2023-01-05T08:17:54.992776Z",
                  "type": "release-sdk-archive",
                  "url": "https://samplelib.com/lib/preview/jpeg/sample-clouds-400x300.jpg",
                  "mime": "image/jpeg",
                  "size": 47419,
                  "originalPath": "sample-clouds-400x300.jpg"
                }
              ],
              "total": 1
            },
            "version": "1.0.0",
            "codeVersion": "1.0.0",
            "contentVersion": "1.0.0",
            "name": "Genesis",
            "description": "Test SDK Release 1.0.0",
            "archive": false
          }
        ]
      }
    },
    "releases": {}
  }
]
```

Example of the `AppV2` metadata, including releases, stored in the database:

```json
{
  "id": "11111111-1111-4111-A111-111111111111",
  "createdAt": "2023-01-01T11:06:19.521Z",
  "owner": {
    "id": "AAAAAAAA-AAAA-4AAA-AAAA-AAAAAAAAAAAA",
    "createdAt": "0001-01-01T00:00:00Z",
    "name": "Hackerman"
  },
  "files": {
    "entities": [
      {
        "id": "aba345d8-de86-4f8c-bfe0-e5a69dd7e7e7",
        "createdAt": "2023-01-01T11:08:34.022464Z",
        "type": "launcher-app-background",
        "url": "https://samplelib.com/lib/preview/jpeg/sample-clouds-400x300.jpg",
        "mime": "image/jpeg",
        "size": 47419,
        "originalPath": "test.jpg"
      },
      {
        "id": "21494b48-c5c4-4b0e-b039-3454fa94e386",
        "createdAt": "2023-01-01T11:08:34.022464Z",
        "type": "launcher-app-card",
        "url": "https://samplelib.com/lib/preview/jpeg/sample-clouds-400x300.jpg",
        "mime": "image/jpeg",
        "size": 47419,
        "originalPath": "test.jpg"
      }
    ]
  },
  "name": "LE7EL",
  "description": "Own your metaverse",
  "external": false,
  "releases": {
    "entities": [
      {
        "id": "827db323-930e-4737-b94e-ff5300669a96",
        "createdAt": "2023-01-04T10:20:44.050Z",
        "files": {
          "entities": [
            {
              "id": "a9a30067-e1e6-4be9-8c3d-b81e1ee3d929",
              "createdAt": "2023-01-23T07:05:52.692874Z",
              "type": "release-archive",
              "url": "https://xxxx..s3-xxxx.amazonaws.com/827db323-930e-4737-b94e-ff5300669a96/a9a30067-e1e6-4be9-8c3d-b81e1ee3d929",
              "mime": "application/zip",
              "size": 1350607527,
              "originalPath": "GameClient-1.0.0-Client-Win64.zip"
            }
          ]
        },
        "version": "1.0.0",
        "codeVersion": "1.0.0",
        "contentVersion": "1.0.0",
        "name": "Genesis",
        "description": "Test App Release 1.0.0",
        "archive": true
      }
    ]
  }
}
```
