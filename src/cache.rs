use std::{
    io::Error,
    path::{Path, PathBuf},
    process,
};

const REPO_URL: &str = "https://github.com/github/gitignore";

fn get_cache_dir() -> Result<PathBuf, Error> {
    // get home directory
    let home = match home::home_dir() {
        Some(path) => path,
        None => panic!("Could not find home directory"),
    };

    let binding = home.join(".github.gitignore");
    let path = Path::new(&binding);

    Ok(path.to_owned())
}

fn is_needs_update() -> Result<bool, Error> {
    let out = process::Command::new("git")
        .arg("rev-list")
        .arg("--count")
        .arg("HEAD..@{u}")
        .output()?;

    return Ok(out.status.success());
}

pub fn prepare_cache() -> Result<(), Error> {
    let path = get_cache_dir()?;
    let path_str = match path.to_str() {
        Some(path) => path,
        None => panic!("Could not convert path to string"),
    };
    if !path.exists() {
        println!("Cloning repository...");
        process::Command::new("git")
            .arg("clone")
            .arg("--depth=2")
            .arg(REPO_URL)
            .arg(path_str)
            .output()?;
    } else if is_needs_update()? {
        println!("Updating repository...");
        process::Command::new("git")
            .arg("-C")
            .arg(path_str)
            .arg("pull")
            .arg("--depth=2")
            .arg("origin")
            .arg("main")
            .output()?;
    }

    println!("Cache updated: {}", path.display());

    Ok(())
}

pub struct LanguageFile {
    pub language: String,
    pub content: String,
    pub file_path: PathBuf,
}

pub fn get_language_file(language: String) -> Result<LanguageFile, Error> {
    let cache_dir = get_cache_dir()?;
    // TODO make this a case-insensitive file search
    let language_file = &cache_dir.join(format!("{}.gitignore", language));
    let content = match std::fs::read_to_string(language_file) {
        Ok(content) => content,
        Err(_) => panic!("Could not read file"),
    };

    Ok(LanguageFile {
        language,
        content,
        file_path: language_file.into(),
    })
}

pub fn get_all_languages_contents(languages: Vec<String>) -> String {
    let mut output = String::new();
    for chosen in languages {
        let language_file = match get_language_file(chosen.to_string()) {
            Ok(info) => info,
            Err(e) => panic!("Error: {}", e),
        };
        let content = language_file.content;
        let sep = "#========================================================================\n";
        output.push_str(format!("\n{sep}# {}\n{sep}\n", language_file.language).as_str());
        output.push_str(&content);
    }
    output = format!("{}\n", output.trim_end());
    output
}
