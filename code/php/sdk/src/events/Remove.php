<?php
/**
 * @copyright Copyright (c) PutYourLightsOn
 */

namespace starfederation\datastar\events;

use starfederation\datastar\Defaults;
use starfederation\datastar\enums\EventType;

class Remove implements EventInterface
{
    use EventTrait;

    public ?string $selector;
    public ?string $settleDuration = null;
    public ?bool $useViewTransition = null;
    public ?array $paths;

    public function __construct(?string $selector = null, array $options = [], array $paths = null)
    {
        $this->selector = $selector;

        foreach ($options as $key => $value) {
            if (property_exists($this, $key)) {
                $this->$key = $value;
            }
        }

        $this->paths = $paths;
    }

    /**
     * @inerhitdoc
     */
    public function getEventType(): EventType
    {
        return EventType::Remove;
    }

    /**
     * @inerhitdoc
     */
    public function getDataLines(): array
    {
        if ($this->selector !== null) {
            $dataLines = ['data: selector ' . $this->selector];

            if ($this->settleDuration !== null && $this->settleDuration !== Defaults::DEFAULT_SETTLE_DURATION) {
                $dataLines[] = 'data: settleDuration ' . $this->settleDuration;
            }

            if ($this->useViewTransition === true) {
                $dataLines[] = 'data: useViewTransition true';
            }

            return $dataLines;
        }

        if ($this->paths !== null) {
            return ['data: paths ' . implode(' ', $this->paths)];
        }

        return [];
    }
}
