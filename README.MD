<img src="assets/logo.png" alt="logo">

---
# <div align="center">KC Checker</div>

<div align="center">Experience the speed of a blazingly fast proxy checker</div>
<div align="center">where checking is done with style </div>

<br>

<div align="center">
<img src="https://img.shields.io/github/downloads/Kuucheen/KC-Checker/total.svg" alt="downloads">
<img src="https://img.shields.io/github/license/Kuucheen/KC-Checker.svg" alt="license">
<img src="https://img.shields.io/github/issues/Kuucheen/KC-Checker.svg" alt="issues">
<br>
<!--<a href="https://discord.gg/7FWAGXzhkC">
  <img src="https://img.shields.io/discord/1196551495571738665?logo=discord&colorB=7289DA" alt="discord">
</a>-->
</div>

<br>

<p align="center">
<img src="assets/preview.gif" alt="preview">
</p>

<br>

# Features

- Multi-Threaded
- Keyword website check
- Verify whether the proxy can access a specific website
- Sleek and intuitive UI for a beautiful user experience
- Highly customizable settings to tailor the tool to your needs

<br>

# Requirements

Your terminal needs to support true color. <br><br>
These are tested terminals on different operating systems:

| System | Tested | Notes |
|--|--|--|
| Windows 11| <div align="center">✅</div> | Works fine |
| Windows 10 | <div align="center">✅</div> | Works without true color
| Linux (ubuntu) | <div align="center">✅</div> | Works fine
| Mac OS | <div align="center">❌</div> | Should work fine

# Getting started

Download [latest release](https://github.com/Kuucheen/KC-Checker/releases/) & unzip the file

**Windows**<br>
Double click on KC-Checker.exe

**Linux**<br>
Open terminal and enter

`chmod +x KC-Checker`

Now run it with

`./KC-Checker`

## Or build it yourself

Install [go](https://go.dev/doc/install) and [git](https://git-scm.com/downloads)

Now open your terminal & clone the project with

`git clone https://github.com/Kuucheen/KC-Checker.git`

**Build**
<details>
  <summary>Using provided scripts</summary>

### Windows

Double click on `start.bat`

### Linux

Open terminal in the directory and make the script executeable

`chmod +x start.sh`

Run it with

`./start.sh`

Now you have an executable. You can run it like described in **Getting Started**
</details>
<details>
  <summary>Using terminal</summary>

Navigate to the directory

    cd KC-Checker

Install dependencies

    go get .

Build with

    go build .

Now you have an executable. You can run it like described in **Getting Started**
</details>

# Customize Settings

Note: Normally you don't have to change these settings

<details>
    <summary>Details for settings.json</summary>

1. **threads**: <br>
   Maximum number of threads

2. **retries**: <br>
   Number of retries for a request

3. **timeout**: <br>
   Timeout duration for requests in ms

4. **keep_alive**: <br>
   If set to true, the request will have the keep-alive header

5. **iplookup**: <br>
   A website that returns the <a href="https://de.wikipedia.org/wiki/Internet_Protocol">ip</a>

6. **judges_threads**: <br>
   Maximum threads for selecting the judge

7. **judges_timeout**: <br>
   Timeout duration of judges in ms

8. **judges**: <br>
   Websites that returns the <a href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers">headers</a> of the request

9. **blacklisted**:<br>
   Websites that contain blacklisted ips. These ips won't be checked

10. **bancheck**: <br>
    If here's a website the program will check if the proxy is able to reach the site. These will land in the `banchecked` directory

11. **keywords**: <br>
    It will check if the website the proxy has opened contains the text given

12. **privacy_mode**: <br>
   If set to true, the proxies will be blured with "*" while checking

13. **copyToClipboard**: <br>
   If set to true, all proxies checked will be copied to clipboard when saving them to a file

</details>

You can blacklist custom proxies by adding them in `blacklisted.txt`

## <p align="center">Other

<div align="center">
<strong>this tool is made for educational purposes only</strong>
<br>
<br>
If you need help with anything join our <a href="https://discord.gg/7FWAGXzhkC">discord server</a>
<br>
<br>
Help me keep my projects open source by <a href="https://ko-fi.com/kuucheen">sponsoring</a> me
</div>
