use std::io;
use std::path::Path;

use super::Format;

pub struct Gif;

impl Format for Gif {
    fn extension(&self) -> &'static str {
        "gif"
    }

    fn read_metadata(&self, path: &Path) -> io::Result<String> {
        Ok(format!("GIF metadata for {}", path.display()))
    }
}
