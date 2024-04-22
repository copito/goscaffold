# GoScaffold

GoScaffold is a command-line tool written in Golang that simplifies project creation by generating boilerplate code from templates, similar to Python's Cookiecutter but leveraging Jinja2 templating engine.

## Installation

To install GoScaffold, you need to have Golang installed on your system. Then you can install it using `go get`:

```bash
go get github.com/copito/goscaffold
```

Moreover, the tool will be compiled to an executable (with the latest version) under the [bin folder](./bin/). So you can use it directly without having to install it if you are testing quickly.

```bash
./bin/goscaffold
```

## Usage

GoScaffold requires a template directory containing your project template files. You can use Jinja2 syntax within these template files to customize the generated code.

To create a new project from a template, use the `goscaffold` command followed by the template directory path and the destination directory:

```bash
goscaffold path/to/template -c path/to/config_file
```

For example:

```bash
goscaffold /home/user/scaffold/example -c example/example.config.yaml
```

This will generate a new project in the ~/Projects/myproject directory using the template located at ~/mytemplate.

## Template Structure

Your template directory should follow a specific structure:

```lua
template/
|-- {{scaffold.project_name}}/
| |-- main.go
| |-- README.md
| |-- ...
|-- hooks
| -- pre_project_gen.sh
| -- post_project_gen.sh
```

Within your template files, you can use Jinja2 variables in double curly braces ({{ }}). These variables will be replaced with user input during project creation.

For example, if your template includes a file named README.md with the content:

```bash
# {{scaffold.project_name}}
```

During project creation, `{{scaffold.project_name}}` will be replaced with the user-provided project name.

## Template Variables

You can define custom variables in a `.scaffold.yaml` file in your template directory.
For example:

```yaml
project_name: "MyProject"
author_name: "John Doe"
email: "john.doe@example.com"
```

These variables can then be used in your template files.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License - see the LICENSE [file](./LICENSE) for details.
Feel free to customize it further based on your project's specifics!
