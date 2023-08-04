#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

vec4 blur(in vec2 uv, in vec2 res) {
    const float pi2 = 6.28318530718; // Pi*2
    
    // GAUSSIAN BLUR SETTINGS {{{
    const float directions = 12.0; // BLUR DIRECTIONS (Default 16.0 - More is better but slower)
    const float quality = 3.4; // BLUR QUALITY (Default 4.0 - More is better but slower)
    const float size = 18.0; // BLUR SIZE (Radius)
    // GAUSSIAN BLUR SETTINGS }}}

    vec2 radius = size / res.xy;
    
    // Pixel colour
    vec4 color = texture(uTexture, uv);
    
    // Blur calculations
    for( float d = 0.0; d < pi2; d += pi2 / directions)
    {
		for(float i = 1.0 / quality; i <= 1.0; i += 1.0 / quality)
        {
			color += texture(uTexture, uv + vec2(cos(d), sin(d)) * radius * i);		
        }
    }
    
    // Output to screen
    color /= quality * directions;
    return color;
}

void main() {
	vec2 uv = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
	fragColor.rgb = blur(uv, uTexBounds.zw).rgb;
	// fragColor.rgb = texture(uTexture, uv, 0.1).rgb;
	fragColor.a = 1.0;
}