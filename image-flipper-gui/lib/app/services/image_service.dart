import 'dart:io';

import 'package:image/image.dart' as img;
import 'package:image_flipper_gui/domain/entities/flip_action.dart';
import 'package:image_flipper_gui/domain/entities/image.dart';
import 'package:image_flipper_gui/domain/extensions/file_extensions.dart';
import 'package:image_flipper_gui/domain/interfaces/i_image_service.dart';
import 'package:path/path.dart' as path;

class ImageService implements IImageService {
  @override
  List<Image> filterImages(String path) {
    final images = <Image>[];

    final directory = Directory(path);

    final files = directory.listSync();

    for (final file in files) {
      if (file is File && file.isImage()) {
        final image = Image.file(file);

        images.add(image);
      }
    }

    return images;
  }

  @override
  Future<void> flipImages(
    List<Image> images,
    FlipAction action,
    String outputDir,
  ) async {
    if (action == FlipAction.original) {
      return;
    }

    for (final image in images) {
      final cmd = img.Command()
        ..decodeImageFile(image.path)
        ..copyFlip(direction: action.flipDirection)
        ..writeToFile(path.join(outputDir, image.name));

      await cmd.executeThread();
    }
  }
}
