# PrivateCloud
Build greatness using your own data

# Core philosophy
- **Access your data** - It's a cloud, d'uh...
- **Sync your data** - Yeah it will have webdav
- **Share your data** - Mom said its not your turn yet with my data
- **Expandable** - That's the juicy part. Low entry barrier + easy plugin development
- **Security** - Yeah we support https like 95% of websites; maybe we'll do 2FA or other stuff in the future

# How to set up PrivateCloud for actual use
1. don't (setup scripts not in development yet; project is still purely developmental)

# How to set up PrivateCloud for development
1. `git clone` or download it using the green button at the top that says "Code.
2. Install go (version 1.23.4 or newer; on Ubuntu: `sudo snap install go --classic`)
3. Open terminal in the projects root directory (where the main.go file is)
4. ```go mod init```
4. ```go mod tidy```
5. ```go run .```

# Troubleshooting
If you're having issues, try `go clean -cache`

# Get in touch
- don't

# Join the team
- don't

# Planned features / core plugins
- Filestorage
- Onlyoffice
- Calendar
- Notes
- Kanban
- Android-App
- Settings panel (User settings & admin settings)
- Permission system (custom permissions? Maybe like in minecraft with \[pluginname].[subset\].\[permission type\])
- Localization (probably a struct that gets initialized with an array of strings mapped to the target language; e.g. log.Printf(localize["Error: Unable to load plugin %s: %v"], path, err))
- Minecraftserver overview
- Maybe a setting to use display the login page as a public dashboard instead? (e.g. Minecraft stats)

THIS WiLL NEVER BE FINNISHED! (bet)
