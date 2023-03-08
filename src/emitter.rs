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
            let final_content = append_lines(original_content, content);
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

/// Append lines from target to source, deduping them, and return the result
/// Does not modify the original source
fn append_lines(target: String, source: String) -> String {
    let original_lines: Vec<&str> = target.lines().into_iter().collect();
    let content_lines: Vec<&str> = source.lines().into_iter().collect();
    let mut filtered_lines = String::new();

    // dedupe
    for line in content_lines.into_iter() {
        if !original_lines.contains(&line) {
            filtered_lines.push_str(&line);
            filtered_lines.push_str("\n");
        }
    }

    // append deduped lines to final output
    let mut final_content = target.clone();
    if !final_content.ends_with("\n") {
        final_content.push_str("\n");
    }
    final_content.push_str(&filtered_lines);
    final_content
}

pub enum EmitStrategy {
    Append,
    Overwrite,
    Skip,
}

#[cfg(test)]
mod tests {

    #[test]
    fn test_append_lines() {
        let original_content = String::from("a\nb\nc");
        let content = String::from("b\nc\nd");

        let expected = String::from("a\nb\nc\nd\n");
        let result = super::append_lines(original_content, content);
        assert_eq!(result, expected)
    }
}
