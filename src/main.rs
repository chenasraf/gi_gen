use analyzer::get_language_candidates;
use cache::get_language_file;
use emitter::emit_file;

use crate::cache::prepare_cache;

mod analyzer;
mod cache;
mod emitter;

fn main() {
    let wd = match std::env::current_dir() {
        Ok(path) => path,
        Err(e) => panic!("Could not get current directory: {}", e),
    };
    match prepare_cache() {
        Ok(_) => println!("Cache prepared"),
        Err(e) => panic!("Error: {}", e),
    }
    let languages = match get_language_candidates(&wd) {
        Ok(result) => result,
        Err(e) => panic!("Error: {}", e),
    };
    let chosen = match languages.get(0) {
        Some(language) => language,
        None => panic!("Could not find language"),
    };
    let content = match get_language_file(chosen.to_string()) {
        Ok(content) => content,
        Err(e) => panic!("Error: {}", e),
    };
    let file = wd.join(".gitignore");
    match emit_file(&file, content, emitter::EmitStrategy::Append) {
        Ok(_) => println!("File emitted: {}", file.display()),
        Err(e) => panic!("Error: {}", e),
    };
}
