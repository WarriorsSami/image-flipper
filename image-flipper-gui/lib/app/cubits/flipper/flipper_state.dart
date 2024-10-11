import 'package:image_flipper_gui/domain/entities/flip_action.dart';
import 'package:image_flipper_gui/domain/entities/image.dart';

sealed class FlipperState {}

final class FlipperInitial implements FlipperState {}

final class FlipperLoadFolderInProgress implements FlipperState {}

final class FlipperLoadFolderSuccess implements FlipperState {
  final List<Image> images;

  FlipperLoadFolderSuccess({required this.images});
}

final class FlipperPreviewFlipImagesSuccess implements FlipperState {
  final List<Image> images;
  final FlipAction action;

  FlipperPreviewFlipImagesSuccess({
    required this.images,
    required this.action,
  });
}

final class FlipperSaveImagesInProgress implements FlipperState {}

final class FlipperSaveImagesSuccess implements FlipperState {
  final List<Image> images;
  final FlipAction action;
  final String outputDir;

  FlipperSaveImagesSuccess({
    required this.images,
    required this.action,
    required this.outputDir,
  });
}

final class FlipperError implements FlipperState {
  final String message;

  FlipperError({required this.message});
}

extension FlipperStateX on FlipperState {
  bool get noFolderSelected =>
      this is FlipperInitial ||
      this is FlipperLoadFolderInProgress ||
      this is FlipperError;

  bool get noFlipActionSelected =>
      this is FlipperInitial ||
      this is FlipperLoadFolderInProgress ||
      this is FlipperLoadFolderSuccess ||
      this is FlipperSaveImagesInProgress ||
      this is FlipperError;
}
