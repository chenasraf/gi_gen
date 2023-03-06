use analyzer::get_language_candidates;

use crate::cache::prepare_cache;

mod analyzer;
mod cache;

fn main() {
    let wd = match std::env::current_dir() {
        Ok(path) => path,
        Err(e) => panic!("Could not get current directory: {}", e),
    };
    match prepare_cache() {
        Ok(_) => println!("Cache prepared"),
        Err(e) => panic!("Error: {}", e),
    }
    match get_language_candidates(wd) {
        Ok(result) => println!("Result: {:?}", result),
        Err(e) => panic!("Error: {}", e),
    }
}
