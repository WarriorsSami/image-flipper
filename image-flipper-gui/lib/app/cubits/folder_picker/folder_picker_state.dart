import 'package:image_flipper_gui/domain/image.dart';

sealed class FolderState {}

final class FolderInitial implements FolderState {}

final class FolderLoading implements FolderState {}

final class FolderLoaded implements FolderState {
  final List<Image> images;

  FolderLoaded({required this.images});
}

final class FolderError implements FolderState {
  final String message;

  FolderError({required this.message});
}
