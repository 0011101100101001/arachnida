mod bmp;
mod gif;
mod jpeg;
mod png;

pub use bmp::Bmp;
pub use gif::Gif;
pub use jpeg::Jpeg;
pub use png::Png;

use std::io;
use std::path::Path;

pub trait Format {
    fn extension(&self) -> &'static str;
    fn read_metadata(&self, path: &Path) -> io::Result<String>;
}
