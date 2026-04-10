<div align="center">
<pre>
░█▀█░█▀▄░█▀█░█▀▀░█░█░█▀█░▀█▀░█▀▄░█▀█
░█▀█░█▀▄░█▀█░█░░░█▀█░█░█░░█░░█░█░█▀█
░▀░▀░▀░▀░▀░▀░▀▀▀░▀░▀░▀░▀░▀▀▀░▀▀░░▀░▀
</pre>
</div>

### [Spider](spider)
> Web Crawler - **Go**

### [Scorpion](scorpion)
> Image Metadata Manager - **Rust**

<hr>

### Bridge
Bridge script that connect these two programs.<br>
First it takes and sends arguments to **Spider** then run **Scorpion** using the
image dowloaded from the directory used in the first one.

#### Usage
```sh
./bridge.py [spider options] URL
```
