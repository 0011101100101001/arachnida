use std::collections::HashMap;
use std::env;
use std::fs;
use std::fs::Metadata;
use std::path::{Path, PathBuf};

pub const DEFAULT: &str = "\x1b[0m";
pub const BOLD: &str = "\x1b[1m";
pub const ITALIC: &str = "\x1b[3m";

pub const RED: &str = "\x1b[31m";
pub const GREEN: &str = "\x1b[32m";
pub const YELLOW: &str = "\x1b[33m";
pub const BLUE: &str = "\x1b[34m";
pub const MAGENTA: &str = "\x1b[35m";
pub const CYAN: &str = "\x1b[36m";
pub const WHITE: &str = "\x1b[37m";

pub struct Image {
    pub name: String,
    pub path: PathBuf,
    pub size: u64,
    pub extension: String,
    pub metadata: Metadata,
}

pub struct Scorpion {
    pub images: HashMap<String, Image>,
}

impl Scorpion {
    pub fn new() -> Self {
        Self {
            images: HashMap::new(),
        }
    }

    pub fn parse_argument(&mut self) -> () {
        for arg in env::args().skip(1) {
            let path = Path::new(&arg);

            let extension = match path.extension().and_then(|s| s.to_str()) {
                Some(ext) => ext.to_lowercase(),
                None => continue,
            };

            if !["png", "jpg", "jpeg", "gif", "bmp"]
                .contains(&extension.as_str())
            {
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
                metadata: meta,
                extension: extension,
            };

            self.images.insert(file_name, image);
        }
    }

    pub fn images(&self) -> &HashMap<String, Image> {
        &self.images
    }

    pub fn print_images(&self) -> () {
        for (name, image) in self.images() {
            println!(
                "{}{}{}: {}{} bytes ({})",
                BOLD, WHITE, name, DEFAULT, image.size, image.extension,
            );
        }
    }
}
