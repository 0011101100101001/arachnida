use crate::scorpion::Scorpion;
use eframe::*;

impl eframe::App for Scorpion {
    fn ui(&mut self, ui: &mut egui::Ui, _frame: &mut Frame) {
        egui::CentralPanel::default().show_inside(ui, |ui| {
            ui.label(format!("{} images loaded", self.images().len()));
            ui.separator();

            egui::ScrollArea::vertical().show(ui, |ui| {
                for (name, image) in self.images() {
                    ui.group(|ui| {
                        ui.label(format!("Name: {}", name));
                        ui.label(format!("Path: {}", image.path.display()));
                        ui.label(format!("Size: {} bytes", image.size));
                        ui.label(format!("Extension: {}", image.extension));
                    });
                    ui.separator();
                }
            });
        });
    }
}
