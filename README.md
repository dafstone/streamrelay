# Streamrelay

Runs an RTMP server that relays video streams to other RTMP servers.

For use with e.g. streaming to multiple streaming services (such as YouTube or Twitch) from the same source stream.

## Usage

Create a file `rtmp-servers.txt` in the *same directory as the `streamrelay` executable, with one `rtmp://` URL per line:

    rtmp://service-1.com/endpoint/key1
    rtmp://service-2.com/endpoint/key2
    rtmp://service-3.com/endpoint/key3

Then start `streamrelay`.

You may optionally pass in the filename as the first argument to the program if you want to use a custom location.