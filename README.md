# Pin

Pin is a TUI base16 theme manager, it applies base16 themes globally using templates.

<img src="./readme-assets/example.gif" alt="Image of pin" />

## Description

Pin is a TUI base16 theme manager heavily inspired by [flavours](https://github.com/Misterio77/flavours), but I wanted something more up to date with the latest [Tinted Theming](https://github.com/tinted-theming) spec. So I built pin as a way to further learn Go.

I have only tested Pin on linux, other OS may have some issues. If you encounter any bugs or have any ideas to improve pin please create an [issue](https://github.com/ClaraSmyth/pin/issues) and let me know, Thanks.

Alternatively feel free to open a [pull request](https://github.com/ClaraSmyth/pin/pulls) if you know how to fix an issue or just want to help out!

## Installation

Install pin with Go

```bash
  go install github.com/ClaraSmyth/pin@main
```

Or download one of the binaries [Here](https://github.com/ClaraSmyth/pin/releases)

## Built Using

- [Go](https://go.dev/)
- [Bubble tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Huh](https://github.com/charmbracelet/huh) 
- [Tinted Theming Schemes](https://github.com/tinted-theming/schemes)

## Acknowledgements

- [flavours](https://github.com/Misterio77/flavours) - Pin is heavily inspired by flavours.
- [Tinted Theming](https://github.com/tinted-theming) - For the base16 colour schemes
- [Charm](https://charm.sh/) - For their great cli tools
- [nap](https://github.com/maaslalani/nap) - I learned a lot from nap about using bubble tea

## Config

Default Config file

```yaml
# Change the default shell and any required args
DefaultShell: sh -c

# Change the default start string to search for when inserting templates
InsertStart: START_PIN_HERE

#Change the default end string to search for when inserting templates
InsertEnd: END_PIN_HERE
```

If you want to change where the config/data is stored by default you can do so with these env variables.

Config files: `PIN_HOME`

Data (cloned schemes etc.): `PIN_DATA`


## Basic Usage/Examples

#### Apps

**You will want to start be creating a new app**

Press **"n"** on the App pane

- **Name** - Name of the application you want to theme
- **Hook** - Any command you want to run after updating the theme
- **Write Method** - Rewrite the entire config file or insert between 2 points
- **Select config file** - Select the apps theme config file

Pressing **"enter"** to select an App will set it to active meaning themes will be applied. If you see an ✗ indicator beside an App it means it is missing either a config file or an active template.

<details>
<summary><b>Insert Write Method?</b></summary>
    <br>
    
If you selected the insert write method you will want to wrap the lines where  you want the template to be inserted in 2 comments. The first must contain the text **"START_PIN_HERE"** and the second must contain **"END_PIN_HERE"**. These insert strings can be changed in the config if required.

**Example**
```
...
# START_PIN_HERE
Any lines you want pin to
overwrite with the template
# END_PIN_HERE
...
```
</details>

---

#### Templates

After creating an app a default Backup template will be created for you, this is just a copy of the config file before being modified.

When creating a new template it will contain the current config file contents or the section between the 2 insert points by default.

You will need to open the template and modify it by adding tags to be able to apply base16 themes to it see the example below for details. You can also find lots of premade templates in the [Tinted Theming](https://github.com/tinted-theming/home) repo. **Templates are .mustache files**.

Pressing **"enter"** to select a template will set it as the active template for the app.

Pressing **"c"** on a template will create a copy of it.

Templates named after a theme will overwrite the active template when that theme is selected. This can be useful for example if you want to have an active template with dynamic tags but overwrite the base16 version of a theme with a hard coded one.
 
<details>
<summary><b>Example Template for Zellij</b></summary>
<br>
    
In this example the **{{base06-hex}}** tag will be replaced with the correct hex color when applying a theme.

```mustache
    themes {
        base16 {
            fg "#{{base06-hex}}"
            bg "#{{base02-hex}}"
            black "#{{base00-hex}}"
            red "#{{base08-hex}}"
            green "#{{base0D-hex}}"
            yellow "#{{base0B-hex}}"
            blue "#{{base0C-hex}}"
            magenta "#{{base0D-hex}}"
            cyan "#{{base0E-hex}}"
            white "#{{base06-hex}}"
            orange "#{{base0A-hex}}"
        }
    }
```
</details>

---

#### Themes 

Pressing **"Alt + p"** on the Themes pane will fetch all the schemes from [Tinted Theming](https://github.com/tinted-theming/home). ***Requires git be  installed**

Pressing **"enter"** will select and apply the theme to the Apps config file using the active template or the template matching the theme name. Once applied the Apps hook will then be run. 

When creating a new theme it will use the current active theme as a base. You should create a new theme when you want to customise it as the schemes pulled from tinted theming will be overwritten anytime you refetch them.

Theme hooks are commands that will be run once per theme. E.g. Change your wallpaper for that theme.

An **✗** indicator means an error occured trying to apply that theme, make sure the theme is formatted correctly.

**Fetching themes will overwrite existing themes!**
If you want to customise a theme you should apply it then create a new theme with it as a base.
Custom themes will not get reset when re-fetching. 

---

#### Cli

You can run also trigger a theme change using the cli by calling pin with a theme name

```bash
pin 'theme name'
```

---

#### Examples

<img src="./readme-assets/example.gif" alt="Image of pin" />

Check out my dotfiles to see how I use pin - [Here](https://github.com/ClaraSmyth/dotfiles)


## License

[MIT](https://github.com/ClaraSmyth/pin/blob/main/LICENSE)


## Things to do

- Improve Error handling
- Clean up code
