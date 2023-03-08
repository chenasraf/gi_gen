use analyzer::get_language_candidates;
use cache::{get_all_languages_contents, prepare_cache};
use emitter::emit_file;

mod analyzer;
mod cache;
mod cli;
mod emitter;

fn main() {
    let args = std::env::args().collect::<Vec<String>>()[1..].to_vec();
    let config = match cli::parse(args) {
        Ok(config) => config,
        Err(e) => panic!("Error: {}", e),
    };
    println!("{}", config);
    let wd = match std::env::current_dir() {
        Ok(path) => path,
        Err(e) => panic!("Could not get current directory: {}", e),
    };
    match prepare_cache() {
        Ok(_) => println!("Cache prepared"),
        Err(e) => panic!("Error: {}", e),
    }
    let languages = match config.languages.len() {
        0 => match get_language_candidates(&wd) {
            Ok(result) => result,
            Err(e) => panic!("Error: {}", e),
        },
        _ => config.languages,
    };
    let output = get_all_languages_contents(languages);
    let file = wd.join(".gitignore");
    match emit_file(&file, output, emitter::EmitStrategy::Append) {
        Ok(_) => println!("File emitted: {}", file.display()),
        Err(e) => panic!("Error: {}", e),
    };
}
