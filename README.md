[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/designed-in-etch-a-sketch.svg)](https://forthebadge.com)

# vimtmpl

Allows you to create template files for scripts.

This is for goofing around... be warned

## Config file

Basic configuration is done in ``~/.template.cfg``


```ini
[DEFAULT]
user = dduck
username = Donald Duck
company = Ducktown
copyright = Donald Duck
mailaddress = d.duck@example.com

license = MIT
mode = 0644

[bash]
description = Bash script
mode = 0755

[bashlib]
extension = .sh
mode = 0644

[go]
extension = .go
mode = 0644

[playbook]
extension = .yml
mode = 0644

[pythonlib]
extension = .py
mode = 0644

[python]
mode = 0755


```

## Templates

Templates are stored in ``~/.templates.d``

Original version:
[vimtmpl-templates](https://github.com/jvzantvoort/vimtmpl-templates.git)

```shell
mkdir -p ~/.templates.d
git clone https://github.com/jvzantvoort/vimtmpl-templates.git defaults
```
