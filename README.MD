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
</div>

<p align="center">
<img src="assets/preview.gif" alt="preview">
</p>

<br>
<br>

-----

### <p align="center">Features</p>

- Multi-Threaded
- Keyword website check
- Verify whether the proxy can access a specific website
- Sleek and intuitive UI for a beautiful user experience
- Highly customizable settings to tailor the tool to your needs

<br>

-----

### <p align="center">Getting started</p>

Install [go](https://go.dev/doc/install)

**Run**
<details>
  <summary>Using provided scripts</summary>
  
  ### Windows
  
  Double click on `start.bat`

  ### Linux

  Open directory in terminal and type `./start.sh`
</details>
<details>
  <summary>Using terminal</summary>

  Navigate to your directory
  > cd your-directory
 
  Install dependencies
  >go get .

  Run
  >go run .
</details>

---

### <p align="center">Customize Settings</p>

Note: Normally you don't have to change these settings

<details>
    <summary>Details</summary>

1. Threads, Retries, Timeout: <br>
    threads: Maximum number of threads<br>
    retries: Number of times to retry a request<br>
    timeout: Timeout duration for requests in ms<br>
   <br>
2. IP Lookup:<br>
    iplookup: A website that returns the [ip](https://de.wikipedia.org/wiki/Internet_Protocol)<br>
<br>
3. Judges:<br>
    judges: Websites that returns the [header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers?retiredLocale=de) of the request<br>
<br>
4. Blacklisted:<br>
    blacklisted: Websites that contain blacklisted ips. These ips won't be checked<br>
<br>
5. Bancheck, Keywords:<br>
    bancheck: If here's a website the program will check if the proxy is able to reach the site. These will land in the `banchecked` directory<br>
    keywords: It will check if the website the proxy has opened contains the text given<br>
</details>

---
### <p align="center">Other</p>

<div align="center">
Help me keep my projects open source by <a href="https://ko-fi.com/kuucheen">sponsoring</a> me
</div>