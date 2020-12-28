# partner

`partner` is a tool for managing coauthors for `git` commits. It allows you to
quickly annotate commits with appropriate `Co-Authored-By` trailers.

## Usage

Add some colleagues to your list of coauthors:

```
# Add a few of my colleagues on GitHub
$ partner manifest gh-add GeorgeMac gavincabbage stuartcarnie

# Add a friend who doesn't use GitHub
$ partner manifest add --id=gemini --email=gemini@strongbeard.org --name="Gemini Strongbeard"

# List all coauthors
$ partner manifest ls
ID            NAME                EMAIL                                          TYPE
gavincabbage  Gavin Cabbage       5225414+gavincabbage@users.noreply.github.com  github
gemini        Gemini Strongbeard  gemini@strongbeard.org                         manual
GeorgeMac     George              1253326+GeorgeMac@users.noreply.github.com     github
stuartcarnie  Stuart Carnie       52852+stuartcarnie@users.noreply.github.com    github
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

`partner` has set a `commit.template` configuration for this repository with
appropriate `Co-Authored-By` trailers. You can view this template like this:

```
$ cat $(git config commit.template)


# Managed by partner
#
# partner-id: gavincabbage
Co-Authored-By: "Gavin Cabbage" <5225414+gavincabbage@users.noreply.github.com>
# partner-id: GeorgeMac
Co-Authored-By: "George" <1253326+GeorgeMac@users.noreply.github.com>
```

The template will be used as a starting point for any commit messages you author
with your party. **Remember:** Using `git commit --message` overrides the entire
commit message and will not use the template.

To clean up:

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
