use crate::style::*;
use std::collections::HashMap;
use std::env;
use std::fs;
use std::fs::Metadata;
use std::path::{Path, PathBuf};
use crate::format::{Bmp, Format, Gif, Jpeg, Png};

pub struct Image<F> {
    pub name: String,
    pub path: PathBuf,
    pub size: u64,
    pub extension: String,
    pub metadata: Metadata,
    pub format: F,
}

type DynImage = Image<Box<dyn Format>>;

pub struct Scorpion {
    pub images: HashMap<String, DynImage>,
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

            let format: Box<dyn Format> = match extension.as_str() {
                "png" => Box::new(Png),
                "jpg" | "jpeg" => Box::new(Jpeg),
                "gif" => Box::new(Gif),
                "bmp" => Box::new(Bmp),
                _ => continue,
            };

            let image = Image {
                name: file_name.clone(),
                path: path.to_path_buf(),
                size: meta.len(),
                metadata: meta,
                extension: extension,
                format: format,
            };

            self.images.insert(file_name, image);
        }
    }

    pub fn images(&self) -> &HashMap<String, DynImage> {
        &self.images
    }

    pub fn print_images(&self) -> () {
        for (_, image) in self.images() {
            _ = image.metadata;
            _ = image.format;
            println!(
                "{}{}{}: {}{} bytes ({})",
                BOLD, WHITE, image.name, DEFAULT, image.size, image.extension,
            );
        }
    }
}
