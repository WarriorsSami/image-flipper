import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:image_flipper_gui/app/cubits/flipper/flipper_cubit.dart';
import 'package:image_flipper_gui/app/cubits/flipper/flipper_state.dart';
import 'package:image_flipper_gui/domain/entities/flip_action.dart';

class ImagesWidget extends StatelessWidget {
  const ImagesWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<FlipperCubit, FlipperState>(
      builder: (context, state) {
        return switch (state) {
          FlipperInitial() => const Center(
              child: Text('No folder selected'),
            ),
          FlipperLoadFolderInProgress() ||
          FlipperSaveImagesInProgress() =>
            const Center(
              child: CircularProgressIndicator(),
            ),
          FlipperLoadFolderSuccess(:final images) => GridView.builder(
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 3,
                crossAxisSpacing: 5,
                mainAxisSpacing: 5,
              ),
              padding: const EdgeInsets.all(8),
              shrinkWrap: true,
              itemCount: images.length,
              itemBuilder: (context, index) {
                final image = images[index];

                return Image.file(
                  File(image.path),
                  fit: BoxFit.contain,
                );
              },
            ),
          FlipperPreviewFlipImagesSuccess(:final images, :final action) ||
          FlipperSaveImagesSuccess(:final images, :final action) =>
            GridView.builder(
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 3,
                crossAxisSpacing: 5,
                mainAxisSpacing: 5,
              ),
              padding: const EdgeInsets.all(8),
              shrinkWrap: true,
              itemCount: images.length,
              itemBuilder: (context, index) {
                final image = images[index];

                return Transform.flip(
                  flipX: action == FlipAction.horizontal ||
                      action == FlipAction.both,
                  flipY: action == FlipAction.vertical ||
                      action == FlipAction.both,
                  child: Image.file(
                    File(image.path),
                    fit: BoxFit.contain,
                  ),
                );
              },
            ),
          FlipperError(:final message) => Center(
              child: Text(
                message,
                style: const TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.bold,
                  color: Colors.red,
                ),
              ),
            ),
        };
      },
    );
  }
}
