use std::{
    error::Error,
    fmt::{self, Formatter},
};

use args::Args;
use getopts::Occur;

#[derive(Debug)]
pub struct Config {
    pub languages: Vec<String>,
    pub auto_discover: bool,
    pub clean_output: bool,
    pub keep_output: bool,
    pub overwrite_file: bool,
    pub append_file: bool,
}

impl fmt::Display for Config {
    fn fmt(&self, f: &mut Formatter) -> fmt::Result {
        write!(
            f,
            "Config{{languages: {:?}, auto_discover: {}, clean_output: {}, keep_output: {}, overwrite_file: {}, append_file: {}}}",
            self.languages, self.auto_discover, self.clean_output, self.keep_output, self.overwrite_file, self.append_file
        )
    }
}

#[derive(Debug, Clone)]
pub enum ConfigErrorType {
    Parse,
    InvalidValue,
    Required,
}

impl fmt::Display for ConfigError {
    fn fmt(&self, f: &mut Formatter) -> fmt::Result {
        write!(
            f,
            "ConfigError{{option: {}, value: {}, message: {:?}, error_type: {:?}}}",
            self.option, self.value, self.message, self.error_type
        )
    }
}

#[derive(Debug, Clone)]
pub struct ConfigError {
    option: String,
    value: String,
    message: Option<String>,
    error_type: ConfigErrorType,
}

impl ConfigError {
    pub fn new(
        option: &str,
        value: Option<&str>,
        error_type: ConfigErrorType,
        message: Option<&str>,
    ) -> ConfigError {
        ConfigError {
            option: option.to_string(),
            value: value.unwrap_or_default().to_string(),
            error_type,
            message: message.map(|s| s.to_string()),
        }
    }
}

pub fn parse(args: Vec<String>) -> Result<Config, ConfigError> {
    let mut parser = create_parser();
    match parser.parse(args) {
        Ok(_) => (),
        Err(err) => {
            return Err(ConfigError::new(
                err.to_string().as_str(),
                None,
                ConfigErrorType::Parse,
                None,
            ))
        }
    };
    let languages = match parser.values_of::<String>("languages") {
        Ok(languages) => languages,
        Err(_) => Vec::new(),
    };
    let auto_discover = if parser.has_value("auto-discover") {
        parser.value_of("auto-discover").unwrap_or(true)
    } else {
        true
    };
    let clean_output = if parser.has_value("clean-output") {
        parser.value_of("clean-output").unwrap_or(true)
    } else {
        false
    };
    let keep_output = if parser.has_value("keep-output") {
        parser.value_of("keep-output").unwrap_or(true)
    } else {
        false
    };
    let overwrite_file = if parser.has_value("overwrite") {
        parser.value_of("overwrite").unwrap_or(true)
    } else {
        false
    };
    let append_file = if parser.has_value("append") {
        parser.value_of("append").unwrap_or(true)
    } else {
        false
    };
    let config = Config {
        languages,
        auto_discover,
        clean_output,
        keep_output,
        overwrite_file,
        append_file,
    };
    Ok(config)
}

fn create_parser() -> Args {
    let mut parser = Args::new(
        "gi_gen",
        "Generate .gitignore files automatically for any project",
    );
    parser.option(
        "l",
        "languages",
        "Comma-separated list of languages to generate .gitignore for",
        "Node,Python,...",
        Occur::Optional,
        None,
    );
    parser.flag(
        "a",
        "auto-discover",
        "Automatically discover languages from project files",
    );
    parser.flag(
        "c",
        "clean-output",
        "Perform cleanup on the output .gitignore file, removing any unused patterns",
    );
    parser.flag(
        "k",
        "keep-output",
        "Do not perform cleanup on the output .gitignore file, keep all the original contents",
    );
    parser.flag(
        "o",
        "overwrite",
        "Overwrite the output .gitignore file if it already exists",
    );
    parser.flag(
        "a",
        "append",
        "Append to the output .gitignore file if it already exists",
    );
    parser.flag(
        "",
        "clear-cache",
        "Clear the local cache of .gitignore files",
    );
    parser.flag("", "all-languages", "List all supported languages");
    parser.flag(
        "",
        "detect-languages",
        "List the automatically-detected languages for the current project",
    );
    parser.flag("h", "help", "Print this help message");

    parser
}
