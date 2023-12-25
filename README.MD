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
        
    cd your-directory
 
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
   Number of times to retry a request<br>
   
3. **timeout**: <br>
   Timeout duration for requests in ms<br>

4. **iplookup**: <br>
    A website that returns the <a href="https://de.wikipedia.org/wiki/Internet_Protocol">ip</a><br>

5. **judges**: <br>
    Websites that returns the <a href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers">headers</a> of the request<br>

6. **blacklisted**:<br>
    Websites that contain blacklisted ips. These ips won't be checked<br>

7. **bancheck**: <br>
   If here's a website the program will check if the proxy is able to reach the site. These will land in the `banchecked` directory<br>

8. **keywords**: <br>
   It will check if the website the proxy has opened contains the text given<br>

</details>

You can blacklist custom proxies by adding them in `blacklisted.txt`


---
### <p align="center">Other</p>

<div align="center">
Help me keep my projects open source by <a href="https://ko-fi.com/kuucheen">sponsoring</a> me
</div>
