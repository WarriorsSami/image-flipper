import 'package:image_flipper_gui/domain/entities/flip_action.dart';
import 'package:image_flipper_gui/domain/entities/image.dart';

abstract interface class IImageService {
  List<Image> filterImages(String path);
  Future<void> flipImages(
      List<Image> images, FlipAction action, String outputDir);
}
