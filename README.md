# partner

`partner` is a tool for managing coauthors for `git` commits. It allows you to
quickly annotate commits with appropriate `Co-Authored-By` trailers.

## Usage

Add some colleagues to your list of coauthors:

```
# Add a few of my colleagues on GitHub
$ partner manifest gh-add GeorgeMac gavincabbage stuartcarnie

# Add a friend who doesn't use GitHub
$ partner manifest add --id=derek --email=derek@strongbeard.org --name="Derek Strongbeard"

# List all coauthors
$ partner manifest ls
ID            NAME               EMAIL                                          TYPE
derek         Derek Strongbeard  derek@strongbeard.org                          manual
gavincabbage  Gavin Cabbage      5225414+gavincabbage@users.noreply.github.com  github
GeorgeMac     George             1253326+GeorgeMac@users.noreply.github.com     github
stuartcarnie  Stuart Carnie      52852+stuartcarnie@users.noreply.github.com    github
```

Activate a few for a pairing session:

```
# Activate George and Gavin
$ partner set GeorgeMac gavincabbage

# List the active coauthors
$ partner status
ID            NAME           EMAIL                                          TYPE
gavincabbage  Gavin Cabbage  5225414+gavincabbage@users.noreply.github.com  github
GeorgeMac     George         1253326+GeorgeMac@users.noreply.github.com     github
```

Make some commits together... ✍️

Clean up:

```
# Unset all coauthors
$ partner clear

# No one is active now
$ partner status
```

## Install

```
$ go get -u github.com/brettbuddin/partner
```

## Environment Variable Overrides

| Environment Variable | Default Value | Description |
| -------------------- | ------------- | ----------- |
| `PARTNER_MANIFEST`   | `~/.config/partner/manifest.json` | Configuration file holding all `add`-ed coauthors. |