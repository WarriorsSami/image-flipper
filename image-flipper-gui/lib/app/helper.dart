import 'dart:io';

import 'package:image_flipper_gui/domain/image.dart';

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

extension FileExtension on File {
  bool isImage() {
    final extension = path.split('.').last;

    return extension == 'jpg' || extension == 'jpeg' || extension == 'png';
  }
}
