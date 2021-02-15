Scrapes T[o]rob products with all details, containing third-party shop's links

# INSTALLATION

Just like other go projects, Clone and run:

```
go build
```

You'll face an executable file in your directory!



# DESCRIPTION

T[o]rob beautifully protects it's pages with Captcha, mainly for fraud detection and avoiding spam clicks (Or at least, that's what we think).

We've created a proxy rotator, which automatically scrapes new proxy from different sources ([Scylla](https://github.com/imWildCat/scylla), for example), and automatically sorts, and validates them, based on proximity of spamminess.

Meanwhile scraping product's sources, a proxy may die, which means, some sources may not be scraped.

There's a repair mode, which identifies corrupt sources and fixes them.

Results of scraping will be saved in a sqlite database.

# OPTIONS

    -workers          Total number of workers (1 is default, and also recommended for small
                      proxy pools)
    -queries          What query to search in torob? only result pages will be scraped. a json
                      file must be passed containing array of queries
    -repair           Repair corrupted sources before scraping new items
    -required-proxies Minimum number of proxies, required in order to start scraping. Program
                      will wait until reaching proxy reuirements
    -h, --help        Face this page in terminal!

