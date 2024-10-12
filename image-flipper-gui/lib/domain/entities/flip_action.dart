import 'package:image/image.dart' as img;

enum FlipAction {
  horizontal,
  vertical,
  both,
  original;

  img.FlipDirection get flipDirection => switch (this) {
        FlipAction.horizontal => img.FlipDirection.horizontal,
        FlipAction.vertical => img.FlipDirection.vertical,
        FlipAction.both => img.FlipDirection.both,
        _ => throw Exception('Invalid flip action'),
      };
}
