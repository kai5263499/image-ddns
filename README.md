# image-ddns
Updates a CloudFlare A record with an IP address in an image. 

I made this after getting a minecraft server for my kids on a hosting provider that charges extra for a static IP address. I'm cheap, and also a programmer, so I wrote this to take the status banner image they provide, OCR it to pull out the text from the image, find the IP address, and use that IP address to update a record in CloudFlare DNS.

The environment variables this command relies on are:   

~~~~bash
export CLOUDFLARE_API_EMAIL="CF_EMAIL"
export CLOUDFLARE_API_KEY="CF_KEY"
export NAME="SUBDOMAIN TO UPDATE"
export ZONE="DOMAIN TO UPDATE"
export IMAGE_URL="IMAGE TO OCR"
export LOG_LEVEL=debug
~~~~

After these are set, you can build and run the image with
~~~~bash
make image-ddns
~~~~
