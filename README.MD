<img src="assets/logo.png" alt="logo">

---
# <div align="center">KC Checker</div>

<div align="center">Experience the speed of a blazingly fast proxy checker</div>
<div align="center">where checking is done with style </div>

<br>

<div align="center">
<!--<img src="https://img.shields.io/github/downloads/Kuucheen/KC-Checker/total.svg" alt="downloads">-->
<img src="https://img.shields.io/github/license/Kuucheen/KC-Checker.svg" alt="license">
<img src="https://img.shields.io/github/issues/Kuucheen/KC-Checker.svg" alt="issues">
<br>
<a href="https://discord.gg/7FWAGXzhkC">
  <img src="https://img.shields.io/discord/1196551495571738665?logo=discord&colorB=7289DA" alt="discord">
</a>
</div>

<br>

<p align="center">
<img src="assets/preview.gif" alt="preview">
</p>

<br>

-----

### <p align="center">Features</p>

- Multi-Threaded
- Keyword website check
- Verify whether the proxy can access a specific website
- Sleek and intuitive UI for a beautiful user experience
- Highly customizable settings to tailor the tool to your needs

<br>

---

### <p align="center">Requirements</p>

Your terminal needs to support true color. <br><br>
These are tested terminals on different operating systems:

| System | Tested | Notes |
|--|--|--|
| Windows 11| <div align="center">✅</div> | Works fine |
| Windows 10 | <div align="center">✅</div> | Works with true color
| Linux (ubuntu) | <div align="center">✅</div> | Works fine
| Mac OS | <div align="center">❌</div> | Should work fine

-----

### <p align="center">Getting started</p>

Download [latest release](https://github.com/Kuucheen/KC-Checker/releases/) & unzip the file

**Windows**<br>
Double click on KC-Checker.exe

**Linux**<br>
Open terminal and enter

`chmod +x KC-Checker`

Now run it with

`./KC-Checker`

------

### Or run it yourself

Install [go](https://go.dev/doc/install) and [git](https://git-scm.com/downloads)

Now open your terminal & clone the project with

`git clone https://github.com/Kuucheen/KC-Checker.git`

**Run**
<details>
  <summary>Using provided scripts</summary>

### Windows

Double click on `start.bat`

### Linux

Open terminal in the directory and make the script executeable

`chmod +x start.sh`

Now you can run it with

`./start.sh`
</details>
<details>
  <summary>Using terminal</summary>

Navigate to the directory

    cd KC-Checker

Install dependencies

    go get .

Run with

    go run .
</details>

---

### <p align="center">Customize Settings</p>

Note: Normally you don't have to change these settings

<details>
    <summary>Details for settings.txt</summary>

1. **threads**: <br>
   Maximum number of threads<br>

2. **retries**: <br>
   Number of retries for a request<br>

3. **timeout**: <br>
   Timeout duration for requests in ms<br>

4. **privacy_mode**: <br>
   If set to true, the proxies will be blured with "*" while checking

5. **iplookup**: <br>
   A website that returns the <a href="https://de.wikipedia.org/wiki/Internet_Protocol">ip</a><br>

6. **judges**: <br>
   Websites that returns the <a href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers">headers</a> of the request<br>

7. **blacklisted**:<br>
   Websites that contain blacklisted ips. These ips won't be checked<br>

8. **bancheck**: <br>
   If here's a website the program will check if the proxy is able to reach the site. These will land in the `banchecked` directory<br>

9. **keywords**: <br>
   It will check if the website the proxy has opened contains the text given<br>

</details>

You can blacklist custom proxies by adding them in `blacklisted.txt`


---
### <p align="center">Other</p>

<div align="center">
If you need help with anything join our <a href="https://discord.gg/7FWAGXzhkC">discord server</a>
<br>
<br>
Help me keep my projects open source by <a href="https://ko-fi.com/kuucheen">sponsoring</a> me
</div>
