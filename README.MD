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

1. **Install Prerequisites:**
   - [Go](https://go.dev/doc/install)
   - [Git](https://git-scm.com/downloads)

2. **Clone the Project:**
   Open your terminal and run:
   `git clone https://github.com/Kuucheen/KC-Checker.git`

3. **(OPTIONAL) Set Up GeoLite2 Database:**
   If you want to determine the country and type (ISP, Datacenter, or Residential), you'll need to download the [GeoLite2 Country and GeoLite2 ASN Database](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) from MaxMind.

   - After downloading, replace the existing databases in the `helper` directory.
   - Ensure they have the same names: `GeoLite2-ASN.mmdb` and `GeoLite2-Country.mmdb`.<br><br>

   > **Note**: These databases are not included in the repository due to their licensing restrictions. However, they are bundled in the [prebuilt version](https://github.com/Kuucheen/KC-Checker/releases/).


<br>**Build**
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

# How to Use

**Navigation**: Use the arrow keys on your keyboard to navigate through the program:

    ↑ Up
    ↓ Down
    ← Left
    → Right
    Press Enter (↵) to select/confirm

**Saving Proxies**: When you save the proxies at the end of a checking progress, they will always be stored in the `output` folder.

# Customize Settings

> [!NOTE]
> For most users, the default settings will probably be sufficient.

<details>
    <summary>Details for settings.json</summary>

1. **threads**: <br>
   Maximum number of threads<br>

2. **retries**: <br>
   Number of retries for a request<br>

3. **timeout**: <br>
   Timeout duration for requests in ms<br>

4. **privacy_mode**: <br>
   If set to true, the proxies will be blured with "*" while checking

5. **autoSelect**: <br>
   When setting one or more protocols are set to true, it will automatically start checking with the selected protocols.

6. **autoSave**: <br>
   **timeBetweenSafes**: Seconds between the autoSave<br>
   **ip:port**: If set to true it will automatically save the proxies with these format every timeBetweenSafes seconds<br>
   **protocol://ip:port**: Same here<br>
   **ip:port;time**: Same here<br>
   **custom**: You can set a custom autoSave format here. If its not empty it will save it with this format. At the end of a checking process theres an option for "CUSTOM", where you can build your own output. You can put the output you built there in here.<br>

7. **timeBetweenRefresh**: <br>
   This is how often the screen updates in ms. If you experience high CPU usage, it is recommended to set it higher.

8. **iplookup**: <br>
   A website that returns the <a href="https://de.wikipedia.org/wiki/Internet_Protocol">ip</a>. This is needed to know if the proxies are Transparent. The website will only be used once when starting the checker.

9. **judges_threads**: <br>
   Maximum threads for selecting the judge

7. **judges_timeout**: <br>
   Timeout duration of judges in ms

8. **judges**: <br>
   These websites are the sites the proxies will connect to. These have to return the <a href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers">headers</a> of the request so the checker can determine what anonymity level the proxy is.

9. **blacklisted**:<br>
   Websites that contain blacklisted ips. These ips won't be checked

10. **bancheck**: <br>
    If here's a website the program will check if the proxy is able to reach the site. These will land in the `banchecked` directory.

11. **keywords**: <br>
    It will check if the website the proxy has opened contains the text given

12. **transport**: <br>
    These are the settings used for making a request. Only touch this option if you know what you are doing.

</details>

You can blacklist custom proxies by adding them in `blacklisted.txt`

## License

The code for this project is licensed under the [MIT License](LICENSE).

The GeoLite2 databases used by this project are provided by MaxMind and are licensed under the [Creative Commons Attribution-ShareAlike 4.0 International License](https://creativecommons.org/licenses/by-sa/4.0/).

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
