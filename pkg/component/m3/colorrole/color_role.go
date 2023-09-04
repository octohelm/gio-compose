package colorrole

type ColorRole string

const (
	Primary            ColorRole = "primary"
	OnPrimary          ColorRole = "on-primary"
	PrimaryContainer   ColorRole = "primary-container"
	OnPrimaryContainer ColorRole = "on-primary-container"

	Secondary            ColorRole = "secondary"
	OnSecondary          ColorRole = "on-secondary"
	SecondaryContainer   ColorRole = "secondary-container"
	OnSecondaryContainer ColorRole = "on-secondary-container"

	Tertiary            ColorRole = "tertiary"
	OnTertiary          ColorRole = "on-tertiary"
	TertiaryContainer   ColorRole = "tertiary-container"
	OnTertiaryContainer ColorRole = "on-tertiary-container"

	Error            ColorRole = "error"
	OnError          ColorRole = "on-error"
	ErrorContainer   ColorRole = "error-container"
	OnErrorContainer ColorRole = "on-error-container"

	SurfaceDim    ColorRole = "surface-dim"
	Surface       ColorRole = "surface"
	SurfaceBright ColorRole = "surface-bright"

	SurfaceContainerLowest  ColorRole = "surface-container-lowest"
	SurfaceContainerLow     ColorRole = "surface-container-low"
	SurfaceContainer        ColorRole = "surface-container"
	SurfaceContainerHigh    ColorRole = "surface-container-high"
	SurfaceContainerHighest ColorRole = "surface-container-highest"

	OnSurface        ColorRole = "on-surface"
	OnSurfaceVariant ColorRole = "on-surface-variant"
	Outline          ColorRole = "outline"
	OutlineVariant   ColorRole = "outline-variant"

	InverseSurface   ColorRole = "inverse-surface"
	InverseOnSurface ColorRole = "inverse-on-surface"
	InversePrimary   ColorRole = "inverse-primary"
)
