use crate::format::{Bmp, Gif, Jpeg, Png};
use crate::style::*;
use std::collections::HashMap;
use std::env;
use std::fs;
use std::fs::Metadata;
use std::io;
use std::path::{Path, PathBuf};

pub enum Format {
    Bmp(Bmp),
    Gif(Gif),
    Jpeg(Jpeg),
    Png(Png),
}

impl Format {
    pub fn read_metadata(&self, path: &Path) -> io::Result<String> {
        match self {
            Format::Png(png) => png.read_metadata(path),
            Format::Jpeg(jpeg) => jpeg.read_metadata(path),
            Format::Gif(gif) => gif.read_metadata(path),
            Format::Bmp(bmp) => bmp.read_metadata(path),
        }
    }
}

pub struct Image {
    pub filename: String,
    pub path: PathBuf,
    pub size: u64,
    pub extension: String,
    pub metadata: Metadata,
    pub format: Format,
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
                    BOLD, RED, arg, RESET
                );
                continue;
            }

            let Ok(meta) = fs::metadata(path) else {
                eprint!(
                    "{}{}Error: {}{} failed to read metadata.",
                    BOLD, RED, arg, RESET
                );
                continue;
            };

            let filename = path
                .file_name()
                .and_then(|s| s.to_str())
                .unwrap_or(&arg)
                .to_string();

            let format = match extension.as_str() {
                "png" => Format::Png(Png::default()),
                "jpg" | "jpeg" => Format::Jpeg(Jpeg::default()),
                "gif" => Format::Gif(Gif::default()),
                "bmp" => Format::Bmp(Bmp::default()),
                _ => continue,
            };

            let image = Image {
                filename: filename.clone(),
                path: path.to_path_buf(),
                size: meta.len(),
                metadata: meta,
                extension: extension,
                format: format,
            };

            match image.format.read_metadata(&image.path) {
                Ok(metdata_string) => println!("{}", metdata_string),
                Err(e) => eprintln!(
                    "Failed to read metadata for {}: {}",
                    image.filename, e
                ),
            }

            self.images.insert(filename, image);
        }
    }

    pub fn images(&self) -> &HashMap<String, Image> {
        &self.images
    }

    pub fn print_images(&self) -> () {
        for (_, image) in self.images() {
            _ = image.metadata;
            _ = image.format;
            println!(
                "{}{}{}: {}{} bytes ({})",
                BOLD,
                WHITE,
                image.filename,
                RESET,
                image.size,
                image.extension,
            );
        }
    }
}
