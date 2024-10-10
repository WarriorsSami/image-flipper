import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:image_flipper_gui/app/cubits/folder_picker/folder_picker_cubit.dart';
import 'package:image_flipper_gui/app/cubits/folder_picker/folder_picker_state.dart';

class ImagesWidget extends StatelessWidget {
  const ImagesWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<FolderCubit, FolderState>(
      builder: (context, state) {
        return switch (state) {
          FolderInitial() => const Center(
              child: Text('No folder selected'),
            ),
          FolderLoading() => const Center(
              child: CircularProgressIndicator(),
            ),
          FolderLoaded(:final images) => GridView.builder(
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
          FolderError(:final message) => Center(
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
