mod format;
mod gui;
mod scorpion;
mod style;

use eframe::NativeOptions;
use scorpion::Scorpion;
use std::env;
use style::*;

fn init_eframe() -> NativeOptions {
    let mut native_options = eframe::NativeOptions::default();

    native_options.viewport = native_options.viewport.with_always_on_top();

    native_options
}

fn main() -> eframe::Result<()> {
    if env::args().len() < 2 {
        eprintln!(
            "{}{}Usage: ./scorpion FILE1 [FILE2 ...]{}",
            BOLD, WHITE, DEFAULT
        );
        std::process::exit(2);
    }

    println!("{}{}{}Scorpion{}", BOLD, ITALIC, MAGENTA, DEFAULT);

    let mut scorpion = Scorpion::new();

    scorpion.parse_argument();
    scorpion.print_images();

    let native_options = init_eframe();

    eframe::run_native(
        "Scorpion",
        native_options,
        Box::new(|_cc| Ok(Box::new(scorpion))),
    )
}
