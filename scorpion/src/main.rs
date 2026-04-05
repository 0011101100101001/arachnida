use std::collections::HashMap;
use std::env;
use std::fs;
use std::path::{Path, PathBuf};

const DEFAULT: &str = "\x1b[0m";
const BOLD: &str = "\x1b[1m";
const ITALIC: &str = "\x1b[3m";

const RED: &str = "\x1b[31m";
const GREEN: &str = "\x1b[32m";
const YELLOW: &str = "\x1b[33m";
const BLUE: &str = "\x1b[34m";
const MAGENTA: &str = "\x1b[35m";
const CYAN: &str = "\x1b[36m";
const WHITE: &str = "\x1b[37m";

struct Image {
    name: String,
    path: PathBuf,
    size: u64,
    extension: String,
}

fn parse_argument() -> HashMap<String, Image> {
    let mut images: HashMap<String, Image> = HashMap::new();

    for arg in env::args().skip(1) {
        let path = Path::new(&arg);

        let extension = match path.extension().and_then(|s| s.to_str()) {
            Some(ext) => ext.to_lowercase(),
            None => continue,
        };

        if !["png", "jpg", "jpeg", "gif", "bmp"].contains(&extension.as_str()) {
            eprint!(
                "{}{}Error: {}{} have not a valid extension.",
                BOLD, RED, arg, DEFAULT
            );
            continue;
        }

        let Ok(meta) = fs::metadata(path) else {
            eprint!(
                "{}{}Error: {}{} failed to read metadata.",
                BOLD, RED, arg, DEFAULT
            );
            continue;
        };

        let file_name = path
            .file_name()
            .and_then(|s| s.to_str())
            .unwrap_or(&arg)
            .to_string();

        let image = Image {
            name: file_name.clone(),
            path: path.to_path_buf(),
            size: meta.len(),
            extension: extension,
        };

        images.insert(file_name, image);
    }

    images
}

fn main() {
    if env::args().len() < 2 {
        eprintln!(
            "{}{}Usage: ./scorpion FILE1 [FILE2 ...]{}",
            BOLD, WHITE, DEFAULT
        );
        std::process::exit(2);
    }

    println!("{}{}{}~Scorpion~{}", BOLD, ITALIC, MAGENTA, DEFAULT);
    let images: HashMap<String, Image> = parse_argument();

    for (name, img) in &images {
        println!(
            "{}{}{}: {}{} bytes ({})",
            BOLD, WHITE,
            name, DEFAULT, img.size, img.extension,
        );
    }
}
