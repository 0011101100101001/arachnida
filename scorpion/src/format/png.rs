use std::io;
use std::path::Path;

use super::Format;
pub struct Png;

impl Format for Png {
    fn extension(&self) -> &'static str {
        "png"
    }

    fn read_metadata(&self, path: &Path) -> io::Result<String> {
        Ok(format!("PNG metadata for {}", path.display()))
    }
}
