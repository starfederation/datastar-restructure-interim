<?php
/**
 * @copyright Copyright (c) PutYourLightsOn
 */

namespace starfederation\datastar\events;

use starfederation\datastar\Consts;
use starfederation\datastar\enums\EventType;

class ExecuteJs implements EventInterface
{
    use EventTrait;

    public string $script;
    public bool $autoRemove = Consts::DEFAULT_EXECUTE_JS_AUTO_REMOVE;
    public string $type = Consts::DEFAULT_EXECUTE_JS_TYPE;

    public function __construct(string $script, array $options = [])
    {
        $this->script = $script;

        foreach ($options as $key => $value) {
            $this->$key = $value;
        }
    }

    /**
     * @inerhitdoc
     */
    public function getEventType(): EventType
    {
        return EventType::ExecuteJs;
    }

    /**
     * @inerhitdoc
     */
    public function getDataLines(): array
    {
        $dataLines = [];

        if ($this->autoRemove !== Consts::DEFAULT_EXECUTE_JS_AUTO_REMOVE) {
            $dataLines[] = $this->getDataLine(Consts::AUTO_REMOVE_DATALINE_LITERAL, $this->getBooleanAsString($this->autoRemove));
        }

        if ($this->type !== Consts::DEFAULT_EXECUTE_JS_TYPE) {
            $dataLines[] = $this->getDataLine(Consts::TYPE_DATALINE_LITERAL, $this->type);
        }

        $dataLines[] = $this->getDataLine(Consts::SCRIPT_DATALINE_LITERAL, $this->script);

        return $dataLines;
    }
}
