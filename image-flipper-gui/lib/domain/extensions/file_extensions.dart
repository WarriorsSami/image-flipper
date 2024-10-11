import 'dart:io';

extension FileExt on File {
  bool isImage() {
    final extension = path.split('.').last;

    return extension == 'jpg' || extension == 'jpeg' || extension == 'png';
  }
}
