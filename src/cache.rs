use std::{
    fs,
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
