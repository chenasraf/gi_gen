use std::{
    fs::File,
    io::{Error, Write},
    path::PathBuf,
};

pub fn emit_file(path: &PathBuf, content: String, strategy: EmitStrategy) -> Result<(), Error> {
    let path_str = match path.to_str() {
        Some(path) => path,
        None => Err(Error::new(
            std::io::ErrorKind::Other,
            "Could not convert path to string",
        ))?,
    };
    let mut file = match File::create(path_str) {
        Ok(file) => file,
        Err(_) => Err(Error::new(
            std::io::ErrorKind::Other,
            "Could not create file",
        ))?,
    };
    match strategy {
        EmitStrategy::Append => {
            let original_content = match std::fs::read_to_string(path_str) {
                Ok(content) => content,
                Err(_) => Err(Error::new(std::io::ErrorKind::Other, "Could not read file"))?,
            };
            let mut final_content = original_content.clone();
            final_content.push_str(&content);
            // TODO smart merge - remove duplicates
            file.write_all(final_content.as_bytes())?;
        }
        EmitStrategy::Overwrite => {
            file.write_all(content.as_bytes())?;
        }
        EmitStrategy::Skip => (),
    }
    Ok(())
}

pub enum EmitStrategy {
    Append,
    Overwrite,
    Skip,
}
